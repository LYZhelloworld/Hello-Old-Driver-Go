using System.Windows;
using System.Net;
using System.IO;
using System.Text.RegularExpressions;
using System.Collections.Generic;
using System.ComponentModel;
using System;
using System.Collections.ObjectModel;
using System.Xml;
using System.ServiceModel.Syndication;

namespace HelloOldDriver
{
    /// <summary>
    /// Interaction logic for MainWindow.xaml
    /// </summary>
    public partial class MainWindow : Window
    {
        public bool IsHTTPS { get; set; } = false;
        public bool IsContentOnly { get; set; } = true;
        public string Domain { get; set; } = "www.liuli.in";
        public int page { get; set; } = 1;
        private BackgroundWorker bw;

        // A list for showing the result
        public ObservableCollection<Result> Result;

        public MainWindow()
        {
            InitializeComponent();
            this.DataContext = this;
            this.Result = new ObservableCollection<Result>();
            ScanResult.ItemsSource = this.Result;
            SetStartStatus(true);
            SetStopStatus(false);
            DisablePagination();
        }

        // Set availability of start button, together with all other controls (except stop button)
        private void SetStartStatus(bool status)
        {
            this.StartButton.IsEnabled = status;
            this.CheckBoxHTTPS.IsEnabled = status;
            this.CheckBoxContentOnly.IsEnabled = status;
            this.TextBoxDomain.IsEnabled = status;
        }

        // Set availability of stop button
        private void SetStopStatus(bool status)
        {
            this.StopButton.IsEnabled = status;
        }

        // Set availability of pagination buttons
        private void SetPagination()
        {
            if(page > 1)
            {
                this.PrevButton.IsEnabled = true;
            }
            else
            {
                this.PrevButton.IsEnabled = false;
            }
            this.NextButton.IsEnabled = true;
        }

        private void DisablePagination()
        {
            this.PrevButton.IsEnabled = false;
            this.NextButton.IsEnabled = false;
        }

        private void StartButton_Click(object sender, RoutedEventArgs e)
        {
            // Disable start button
            SetStartStatus(false);
            // Disable pagination
            DisablePagination();
            // Clear the list
            Result.Clear();
            var s = new RssScanner(Domain, IsHTTPS ? "https" : "http", page);

            bw = new BackgroundWorker
            {
                // Report progress when finished processing one page
                WorkerReportsProgress = true,
                // Supports cancellation when clicking stop button
                WorkerSupportsCancellation = true
            };

            // Read RSS
            List<string> pages = s.Scan();
            // Update progress bar
            this.ScanProgress.Minimum = 0;
            this.ScanProgress.Maximum = pages.Count;
            this.ScanProgress.Value = 0;

            bw.DoWork += new DoWorkEventHandler(delegate(object o, DoWorkEventArgs args)
            {
                BackgroundWorker b = o as BackgroundWorker;
                int progress = 0;

                foreach (string pageLink in pages)
                {
                    // Check if the user clicked "Stop"
                    if (b.CancellationPending)
                    {
                        args.Cancel = true;
                        return;
                    }
                    // If the user clicked Stop after this line, it will be processed at next turn of loop

                    // Get page content
                    string pageContent = s.ScanPage(pageLink);
                    // Get title and magnet links
                    var a = new Analyzer(pageContent, IsContentOnly);
                    var title = a.GetTitle();
                    var magnets = a.GetMagnetLinks();

                    if (magnets.Count > 0)
                    {
                        // This page has magnet links
                        // Do in UI thread
                        App.Current.Dispatcher.Invoke((Action)delegate
                        {
                            // Links = all links separated with newlines
                            Result.Add(new Result() { Title = title, Links = string.Join(Environment.NewLine, magnets) });
                        });
                    }
                    // Update progress
                    progress++;
                    b.ReportProgress(progress);
                }
            });

            bw.ProgressChanged += new ProgressChangedEventHandler(delegate(object o, ProgressChangedEventArgs args)
            {
                // Change the value of progress bar
                App.Current.Dispatcher.Invoke(() =>
                {
                    this.ScanProgress.Value = args.ProgressPercentage;
                });
            });

            bw.RunWorkerCompleted += new RunWorkerCompletedEventHandler(delegate (object o, RunWorkerCompletedEventArgs args)
            {
                // Disable stop button
                SetStopStatus(false);
                // Enable start button
                SetStartStatus(true);
                // Push the progress bar to 100% (even if it is cancelled)
                App.Current.Dispatcher.Invoke(() =>
                {
                    this.ScanProgress.Value = this.ScanProgress.Maximum;
                    SetPagination();
                });
            });

            bw.RunWorkerAsync();
            // Enable stop button after the async work starts
            SetStopStatus(true);
        }

        private void StopButton_Click(object sender, RoutedEventArgs e)
        {
            // Disable stop button before cancelling
            SetStopStatus(false);
            bw.CancelAsync();
        }

        private void PrevButton_Click(object sender, RoutedEventArgs e)
        {
            if(page > 1)
            {
                page--;
            }
            StartButton_Click(sender, e);
        }

        private void NextButton_Click(object sender, RoutedEventArgs e)
        {
            page++;
            StartButton_Click(sender, e);
        }

        private void ListViewItem_MouseDoubleClick(object sender, System.Windows.Input.MouseButtonEventArgs e)
        {
            int index = ScanResult.Items.IndexOf((sender as FrameworkElement).DataContext);
            Clipboard.SetText(Result[index].Links);
            MessageBox.Show(Application.Current.Resources["CopiedToClipboard"] as string + Environment.NewLine + Result[index].Links, Application.Current.Resources["AppName"] as string, MessageBoxButton.OK, MessageBoxImage.Information);
        }
    }

    [Obsolete("Use RssScanner instead.")]
    class Scanner
    {
        private readonly string url;

        public Scanner(string domain, string protocol)
        {
            // Protocol should be "http" or "https"
            // Url will be <protocol>://<domain>/wp/{0}.html
            this.url = string.Format("{0}://{1}/wp/", protocol, domain) + "{0}.html";
        }

        public string Scan(int id)
        {
            // Put the article ID into the placeholder
            string fullUrl = string.Format(url, id);
            WebClient client = new WebClient();
            try
            {
                Stream data = client.OpenRead(fullUrl);
                StreamReader reader = new StreamReader(data);
                string result = reader.ReadToEnd();
                data.Close();
                reader.Close();
                return result;
            }
            catch(WebException)
            {
                return "";
            }
        }
    }

    class Analyzer
    {
        private readonly string resultText;
        private readonly bool isContentOnly;
        private static readonly string MAGNET_PREFIX = "magnet:?xt=urn:btih:";

        public Analyzer(string resultText, bool isContentOnly)
        {
            this.resultText = resultText;
            this.isContentOnly = isContentOnly;
        }

        public List<string> GetMagnetLinks()
        {
            // Processed result text
            string newResultText = this.resultText;
            if(this.isContentOnly)
            {
                // Find the content part
                // There are false strings that are not magnet links but do not appear in the content part
                // <div class="entry-content">...<div class="same_cat_posts">
                Regex rxContent = new Regex(@"\<div(?:\s+?)class=""entry-content""(?:\s*?)\>(?<content>.*?)\<div(?:\s+?)class=""same_cat_posts""(?:\s*?)\>", RegexOptions.Singleline);
                var match = rxContent.Match(this.resultText);
                if(match.Success)
                {
                    newResultText = match.Groups["content"].Value;
                }
                else
                {
                    newResultText = "";
                }
            }
            // Find 40-character magnet links
            Regex rxMagnetLinks = new Regex(@"[0-9a-fA-F]{40}");
            var matches = rxMagnetLinks.Matches(newResultText);
            var result = new List<string>();
            for(int i = 0; i < matches.Count; i++)
            {
                // magnet:?xt=urn:btih:<magnet_link_in_lower_case>
                result.Add(MAGNET_PREFIX + matches[i].Value.ToLower());
            }
            return result;
        }

        public string GetTitle()
        {
            // Find <title>...</title>
            Regex rxTitle = new Regex(@"\<title\>(?<title>.+?)\</title\>");
            var match = rxTitle.Match(this.resultText);
            if(match.Success)
            {
                return match.Groups["title"].Value;
            }
            return "";
        }
    }

    public class Result
    {
        public string Title { get; set; }
        public string Links { get; set; }
    }

    class RssScanner
    {
        private readonly string url;

        public RssScanner(string domain, string protocol, int page = 1)
        {
            var url = string.Format("{0}://{1}/wp/?feed=rss", protocol, domain);
            if(page > 1)
            {
                url += string.Format("&paged={0}", page);
            }
            this.url = url;
        }

        public List<string> Scan()
        {
            try
            {
                XmlReader reader = XmlReader.Create(this.url);
                SyndicationFeed feed = SyndicationFeed.Load(reader);
                reader.Close();
                List<string> result = new List<string>();
                foreach (SyndicationItem item in feed.Items)
                {
                    result.Add(item.Links[0].Uri.ToString());
                }
                return result;
            }
            catch (FileNotFoundException)
            {
                return new List<string>();
            }
        }

        public string ScanPage(string fullUrl)
        {
            WebClient client = new WebClient();
            try
            {
                Stream data = client.OpenRead(fullUrl);
                StreamReader reader = new StreamReader(data);
                string result = reader.ReadToEnd();
                data.Close();
                reader.Close();
                return result;
            }
            catch (WebException)
            {
                return "";
            }
        }
    }
}
